[![ci](https://github.com/OkieOth/gordered-map/actions/workflows/test.yml/badge.svg?branch=main&event=push)](https://github.com/OkieOth/gordered-map/actions/workflows/test.yml)
[![go report card](https://goreportcard.com/badge/github.com/OkieOth/gordered-map)](https://goreportcard.com/report/github.com/OkieOth/gordered-map)


# gordered-map
A simple wrapper to provide an ordered golang type to have fixed order map.

This package doesn't pretend to be performent nor is the usage convinient, but it
helps me in use-cases where I have to read a JSON file to a generic dictionary,
manipulate the content and want to have after serialization the initial order of
the input.
In other words in tries to tackle the fact that golang maps only have a random
order when you travers them ... no I don't complain X-D
