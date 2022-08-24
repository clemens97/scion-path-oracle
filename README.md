# SCION Path Oracle

SCION Path Oracle is an HTTP service deriving dynamic path metrics in [SCION](https://scion-architecture.net/).
networks.
Applications can donate their experienced performance on a path to the Path Oracle, which derives an internal network
view from those donations. Other clients may query the Path Oracle for scores to a given destination AS to perform a
more sophisticated path selection.

![System Overview](doc/overview.drawio.svg)

This service gets implemented in the course of a master thesis at
the [Otto von Guericke University](https://www.ovgu.de/en/)
in Magdeburg.

## Scoring Services

- [x] throughput
- [ ] latency
- [ ] loss
- [ ] ...