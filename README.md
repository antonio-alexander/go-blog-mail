# go-blog-mail (github.com/antonio-alexander/go-blog-mail)

This is a weird one. I've got this project and I wanted a way to verify code that would send emails; in context, I wasn't super interested in _receiving_ emails, but in order to confirm that the code sending the email was successful you have to be able to receive it. Furthermore, this solution has to be closed loop it shouldn't have any external dependencies.

Yes...it's totally possible to use external dependencies to confirm this, things like amazon sns or simply gmail; but those aren't great for a lot of practical reasons:

- without a closed solution, you'll have to worry about how to hide the secrets (username/password)
- you don't have to worry about cleaning up emails after testing
- you have more control over the configuration of the mail server

## References

- [https://pkg.go.dev/net/smtp](https://pkg.go.dev/net/smtp)
- [https://opensource.com/article/21/5/alpine-linux-email](https://opensource.com/article/21/5/alpine-linux-email)
- [https://docker-mailserver.github.io/](https://docker-mailserver.github.io/)
- [https://stackoverflow.com/questions/75247361/error-while-running-docker-mailserver-via-docker-compose-up-you-need-at-least-o](https://stackoverflow.com/questions/75247361/error-while-running-docker-mailserver-via-docker-compose-up-you-need-at-least-o)
- [https://www.ifthenel.se/self-hosted-mail-server/](https://www.ifthenel.se/self-hosted-mail-server/)
- [https://www.iops.tech/blog/postfix-in-alpine-docker-container/](https://www.iops.tech/blog/postfix-in-alpine-docker-container/)
- [https://github.com/jeboehm/docker-mailserver](https://github.com/jeboehm/docker-mailserver)
- [https://mailu.io/2.0/](https://mailu.io/2.0/)
- [https://quay.io/repository/instrumentisto/dovecot](https://quay.io/repository/instrumentisto/dovecot)
- [https://tvi.al/simple-mail-server-with-docker/](https://tvi.al/simple-mail-server-with-docker/)

## Trial and Error

Originally, I wanted to solve this problem in a way that was really easy to _see_, I wanted to have a customer docker image starting from alpine/busybox and then add all of the dependencies for smtp, pop3, imap etc., and then have custom configuration for each one. I got halfway there, but what I found (which in hindsight was obvious) was that...it's __difficult__ to configure a mailserver; even one without security that you just bring up and teardown, its _hard_.

So In the end, I settled on using [ghcr.io/docker-mailserver/docker-mailserver](https://github.com/docker-mailserver/docker-mailserver/pkgs/container/docker-mailserver); because it not only did most of this out of the box and used Docker, but I was able to get something semi-functional after 10m of reading through the documentation. By using this solution, I was able to focus on the things that matter more: determining which library/api would provide me the functionality to talk in each of these protocols using Go.

Sometimes, I think there's value in 100% owning your solution (even for testing); but sometimes...in the interest of time, it's better that you get something that works and lets you focus.

## Protocols for Sending and Receiving Emails

Before we get into the technical portions of how to do this in Go, I think it's good to start with the different mail protocols and what they're expected to do; and this isn't because I'm trying to educate you, but understanding what the protocols are makes it easier to conceptualize the pieces you need in order to test something. I'll begin this with the following:

> Mail protocolS are difficult and they are legion. I imagine you already know this, but understand that it's not a common thing to have different protocols to read and write the same thing; there's not just more than one protocol to handle electronic mail, but there are different protocols who's job is simply to receive or simply to send

With that said, here are the common mail protocols that we want to work with:

- pop3: this is the "post office protocol" version 3 which is used _only_ to read/retrieve email
- imap: this is the "internet message access protocol" which is __ALSO__ used _only_ to read/retrieve email
- smtp: this is the "simple mail transfer protocol" which can be used to __BOTH__ send/receive email although it's often only used to send email
- mapi: this is the "messaging application programming interface" this is also an email protocol; that can send/receive emails, but it's used almost exclusively to talk to exchange; we're not going to really chat about it beyond letting you know that it _exists_

Our general goal is trying to develop the sandbox infrastructure (i.e., infrastructure as code) to be able to verify an application's ability to send and receive email.

## Choosing a Solution

I did a little research about the protocol itself and although I think it's possible to write a Go client against pop3/smtp as those protocols are well documented and generally open source, a lot of other people have already done that work; hell Go has some existing first party support for email protocols (generally auxiliary functionality rather than direct support).

With that said, when I'm vetting an existing library, I generally look to the following to guide my hand:

1. Does it work?
2. Is it open source?
3. How painful is integration?
4. Does it have relatively recent support?
5. Is it wrapping another package/library?

When trying to choose a library, I priorize functionality over everything, it's more important that it works and it helps if it works within the first 10m (I have a short attention span). Although most recent packages are open source software (OSS), if it isn't I won't use it; using non-OSS software becomes a pain if you ever need to license the software beyond it being a proof of concept; if your software isn't making money or being used by someone trying to not pay the people who have a license, no-one really cares.

Finally, I care about long-term support, does it have any, if we find a bug who can we trust to fix it (or at least someone to rage at)? It doesn't need a lot of support and at times, it's even OK if it's an archived library since we can go vendor, fork etc (unless there's a security issue, most packages are safe to be unsupportable). It's important that the way in which it's integrated is relatively closed, packages that don't force you to implement [leaky abstractions](https://en.wikipedia.org/wiki/Leaky_abstraction) are generally really good, in addition packages that don't force you to create a pointer for everything (although this is an ok implementation); is also a hallmark of a _good_ package.

I've chosen the following client packages to integrate with (all of these have an MIT license):

- [github.com/taknb2nch/go-pop3](https://github.com/taknb2nch/go-pop3): for pop3; this is a REALLY old package that has a TON of forks; this clearly doesn't have any current support (it doesn't even have a go module); but it's been forked plenty of times and worked out the box; I'm sure there are other pop3 packages, but it worked and integration was simple. MIT license
- [github.com/emersion/go-imap](https://github.com/emersion/go-imap): this simply worked; I can't stress that enough; it has current support
- [github.com/go-mail/mail](https://github.com/go-mail/mail): this was primarily for smtp, and is a fork of gomail; this fork has "current" support (since 2018) and it's an OK starting point; if a bug fix needed to be implemented, I'd consider forking gomail and trying to figure out what they fixed and then applying my fix (or submitting a PR)

I think it's important to note that these are _client_ packages, in order for this to work we don't just need a client, we need a server too. There are a handful of "server" solutions in Go, with varying levels of functionality. What I found was that the _configuration_ process of these servers was crazy. I have some system admin experience, but considering going through the trouble of configuration brought up a lot of bad memories. So in searching I prioritized something simple, something I didn't have to homebrew and preferably, something that either had Docker support or I could Dockerize.

Originally, I was just going to use busybox, and install everything from scratch (can't be that hard); it was that hard. So I settled on an out of the box solution: [github.com/docker-mailserver/docker-mailserver](https://github.com/docker-mailserver/docker-mailserver). And this worked out of the box with minimal effort; I didn't have to pore over configuration and I was able to flip a few switches to get support for pop3, imap and smtp all in one place. And more importantly, I could start this environment from nothing; so it was perfect for my sandbox environment.

> I'll preface this with most people who want a production-ready Dockerized mail solution are f**cking crazy. I would never, but I'm happy that someone did it and it supports my need for a sandbox environment

## Getting the Server up and running

```sh
docker exec -ti mailserver setup email add user@example.com
```

## Implementing a Solution

## Security Considerations
