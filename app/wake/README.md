# Wake
`Wake` is a command-line tool that listens to `PillarMQ` (a RabbitMQ exchange) for various `event`s coming out of `Pillar`.

For each event type, it then decides how to react. Each reaction takes place in its own `goroutine`.

