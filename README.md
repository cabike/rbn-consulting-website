# RBN Consulting

An Astro static site with Tailwind CSS and React islands. The contact endpoint is a separate Go HTTP serverless handler.

## Development

Requires Node.js 22.23.4+ (or a supported newer release) and Go 1.24+.

```sh
npm ci
npm run dev
npm run build
go test ./...
```

`npm run build` generates the static website in `dist/`. Astro's development and preview servers do not run the Go handler.

## Contact form deployment

Deploy `api/contact.go` as a Go HTTP serverless function and route `POST /api/contact` on the same origin as the static site to its exported `Handler`. Configure the function with `RESEND_API_KEY`, `CONTACT_FROM_EMAIL` (a sender on a verified Resend domain), and `CONTACT_TO_EMAIL`. The handler sends inquiries through the Resend email API; it returns an error if delivery is not configured or fails. Keep these values in the hosting provider's secret configuration, not in the repository.

The form accepts only basic contact details; visitors should not submit sensitive information. Configure rate limiting and abuse protection at the hosting edge before exposing the endpoint publicly.
