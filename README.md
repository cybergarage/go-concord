# Concord

![](https://img.shields.io/badge/status-Work%20In%20Progress-8A2BE2)
![GitHub tag (latest SemVer)](https://img.shields.io/github/v/tag/cybergarage/Concord)
[![test](https://github.com/cybergarage/Concord/actions/workflows/make.yml/badge.svg)](https://github.com/cybergarage/Concord/actions/workflows/make.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/cybergarage/Concord.svg)](https://pkg.go.dev/github.com/cybergarage/Concord)
[![codecov](https://codecov.io/gh/cybergarage/Concord/graph/badge.svg?token=GOLCBMUVB1)](https://codecov.io/gh/cybergarage/Concord)


`Concord` is a coordination service designed to simplify metadata management and distributed system synchronization.
In large-scale distributed environments, multiple nodes must operate cooperatively, requiring reliable mechanisms for agreement, synchronization, and state management. Coordination services address this need by centralizing essential metadata and providing consistent operations.

**Note:** 🌱 This is a spare-time hobby project, so progress may be slow and changes may appear irregular. Thank you for your patience 🙂

## Key Features

Concord provides:

# Architecture & Concepts

Core architecture and design documents:
- [Design Concepts](doc/concept.md)
  - [Data Model](doc/data-model.md)
  - [Storage Concept](doc/storage-concept.md)
  - [Consistency Model](doc/consistency-model.md)
  - [Coordinator Concept](doc/coordinator-concept.md)
