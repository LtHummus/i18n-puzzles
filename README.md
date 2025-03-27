# i18n-Puzzles

This repo contains my implementations to solve the puzzles at https://i18n-puzzles.com. My self-imposed goal was to only use the built in Go standard libraries (plus the expanded ones hosted at golang.org). The only exception to this is https://github.com/stretchr/testify for writing unit tests.

DISCLAIMER: This isn't the cleanest code, but it does solve the problem. At some point, I should go back and clean some stuff up...but not right now :)

## Input Downloader

This repo also contains my downloader for automatically downloading an importing the input data from the site. The downloader will look for your session authentication token in `~/.i18n-puzzles/.token`. You'll have to get this yourself (probably from your browser's tools). Inputs -- since they never change -- will also be cached locally in that `~/.i18n-puzzles` directory. Contact with the server will only be made for the first time you load that puzzle's input data. You can set the flag
`input.RealData` or `input.TestData` to get the real or test data as required.
