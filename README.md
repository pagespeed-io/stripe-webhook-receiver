# Stripe Webhook Receiver

This little webserver listens for incoming [Stripe] webhooks and sends
out [Pushover] notifications to the team, based on the type of data.

[More about Stripe webhooks and their payload format.][stripe-doc]

  [Stripe]:     https://stripe.com
  [Pushover]:   https://pushover.net
  [stripe-doc]: https://stripe.com/docs/webhooks
