# aspenjames.dev -- Static Content

This directory contains the static website content for aspenjames.dev. See the
[main project README][gh] for more details.

## Structure

HTML templates are in `/content`, layouts are in `/content/layouts`. CSS files
are in `/css`, which are built by tailwind into `/static`. Only the `/content`
and `/static` directories will be deployed to the storage bucket.

## Alpine

Alpine is "installed" via CDN in the `layouts/main.html` template. This is
also where the Alpine stores are initialized. Store properties prefixed with
an underscore (`_`) are "private" -- this is a convention and not a rule. A
store's public API is any property not prefixed with an underscore, except for
the special case `init()`.

## Tailwind

Tailwind is configured at `tailwind.config.js` & the input CSS is `input.css`.
These drive the build process, which produces `site.css`. The built CSS file
is not checked in to GitHub, since this is a build artifact. Build with `npm
run build-css`.

[gh]: https://github.com/aspenjames/aspenjames.dev
