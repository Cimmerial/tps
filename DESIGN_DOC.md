
# TPS

## Goal
- TPS is going to be a central hub where I can add/create projects then run commands from them simply, without needing all the jargen.

## Features
- It will be a CLI, where I start it up, then can nav to projects, and within each project I have a list and or directory of commands which i can run, for example in one of my projects I often run locally, instead of needing to type out `npm dev run:mobile` I can just use command `run` or `r`.
- It will be able to open, close, and manage other terminals, where it can run sequences of commands and wait in between for certain events to happen.
- Projects and commands should be exportable and importable.

## Language
- Was thinking C++ or py but perhaps rust

## Questions
- How to make it so i can run cmd from any terminal on computer and it will start up and open?
- Can I tie each project in it to some sort of name/token where if project is moved we can still find it?
- I use a lot of Xcode, can we make it where we use the CLI to run xcode without needing to interact with it directly?
- Can we make it so this works from computer to computer opposed to only on this device so if i get a new one it still does what it should?