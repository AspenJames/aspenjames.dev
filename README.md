# aspenjames.dev

Personal site for Aspen James -- software developer, tree hugger.

## About

I enjoy having a small corner of the internet to express myself and explore
new concepts and technologies. I've "re-written" this website many times as
I've progressed and grown as an engineer. This time around, I want to build it
in such a way as to intentionally highlight some of the things I find most
important in software delivered over the web -- clean, minimal, fast,
available. The goals may change and evolve as I go, but these are the guiding
factors behind the decisions I've made so far:

- [Fastly Compute@Edge][compute]
  - Extremely fast & available edge network
  - Efficient caching & delivery
- [Go][go]
  - Fast, enjoyable language, good tooling & support
  - Compiles to static binary
  - [Fiber][fiber] framework is great
- Static HTML templates
  - Speed, simplicity
  - Go HTML templates widely supported
- Markdown for blog content
  - Easy to write content in md
  - Easy to build pipeline to transform content -> HTML
- Minimal JavaScript
  - [Alpine.js][alpine] - super minimal
  - Necessary for some interactivity & polish
- [Tailwind CSS][tailwind]
  - Utility-based CSS modules
  - Minimal final bundle size

The built site will ultimately live as a set of static files (templates) in a
storage bucket. A build & deploy pipeline could be configured to:

- Compile & minify CSS
  - Utilize Tailwind tooling
- Upload to storage
- Issue cache invalidation

Rendering of templates will be done by the Go Fiber server -- this will be a
simple & fast file server in front of the static files and templates in the
storage bucket. The Compute@Edge will handle caching & delivery, delegating to
the Go server when necessary.

## Structure

The source will be organized like this:

```
.
├── content               -> Website files
│  ├── css                -> Tailwind input files
│  ├── static             -> Static files
│  ├── tailwind.config.js -> TailwindCSS config
│  └── templates          -> Go HTML Templates
├── infra                 -> Infrastructure code
├── main.go               -> Go template renderer / file server
└── src                   -> Compute@Edge project
   └── edge
```

## Go API

The Go API/file server can be run locally with `fiber dev` or with `docker`.
Flags may be supplied to configure values for the application:

* `content` -- filepath to directory containing website content; may be absolute or relative
* `domain` -- website domain; used to set cookie
* `port` -- port on which the server will listen

[alpine]: https://alpinejs.dev/
[compute]: https://developer.fastly.com/learning/compute
[comrak]: https://docs.rs/comrak/latest/comrak/index.html
[fiber]: https://gofiber.io/
[go]: https://go.dev/
[tailwind]: https://tailwindcss.com/
[tera]: https://tera.netlify.app/
