# nubes

_nubes is cloud in Latin_</small>

Holistic self hosting

## What is self hosting?
The self hosting movement is made of people that try to move away from the traditional content silos (mainly Google, Apple, Microsoft) and to retain ownership of our data.

### What kind of data are you talking about?

Photos, documents, emails, music, movies, ebooks, your financial plans, your fitness progress, your passwords to other services.

### Why would you do this?

- **Privacy**: all the data you store in "the Cloud" is continuously monitored and analyzed. Your data is extracted and sold to the best bidder. Does Google need to know how much you keep in your bank account?
- **True ownership**: do you actually own that movie you bought on Google Play Movies? Not really. If you lose access to your account, or Google bots decide you are not worthy, all your purchases are lost. 
  _[This](https://twitter.com/TheyWhoRemain/status/1368955992105484290) [happens](https://twitter.com/miguelytob/status/1315749803041619981) [all](https://www.youtube.com/watch?v=pWaz7ofl5wQ) [the](https://twitter.com/Demilogic/status/1358661840402845696) [time](https://news.ycombinator.com/item?id=25899814)_
- **No lock in**: if you had a Google Pixel and tried to move all your data to an iPhone or vice-versa, you know the pain. These services try their best to hold onto your data and make it hard to leave. GDPR helps, if you live in Europe.
- **Legacy**: what is not owned cannot be given. Not to friends, not to heirs.

[Learn more on Reddit](https://reddit.com/r/selfhosted)

## Self hosting is hard
You have to buy your hardware, keep it running, replace it when it breaks. And when it's a disk breaking, have a plan for replacing it without losing your data.

You have to pick the software, configure it, make it work with each other.

You have to monitor it, secure it, troubleshoot it.

But hey, it's also really cool.

## So what is this?
This project (temporarily named _nubes_) tries to bring a holistic approach to self hosting. 

Having to configure and setup 30 different applications might be fun. Making them work togheter is less fun. Using 30 different applications with no clear, uniform design is annyoing.

Your data is more valuable when it works together. Nubes helps you to manage _all your data_. Photos, music, passwords, meal history, money.

## Goals

### Microservices-like

The first idea was to build one majestic monolithic app, in one language (I love Ruby, but maybe Rust, or Crystal).

However, even if the whole experience of the system should be of a homogeneous single application, it doesn't mean it can't be built of different smaller applications.

This comes a few advantages:
- You don't have to run it all (I don't care for movies, don't waste CPU and memory on it)
- You can use different programming languages where it makes sense.
- A part of the system can more easily be replaced by another, if it follows the API contract.

### Easy to configure

No one should be spending days to link everything together. Assuming you already have the hardware configured, it should be as easy as:

```bash
docker-compose up -d
```

followed by a web based wizard.

### Easy to backup

The various apps shall exclusively be writing to a single shared file/object repository (directory). Saving that directory elsewhere should be enough to fully backup the system.
