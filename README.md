# RPC
[![GoDoc](https://godoc.org/gopkg.in/jcmturner/rpc.v0?status.svg)](https://godoc.org/gopkg.in/jcmturner/rpc.v0) [![Go Report Card](https://goreportcard.com/badge/gopkg.in/jcmturner/rpc.v0)](https://goreportcard.com/report/gopkg.in/jcmturner/rpc.v0) [![Build Status](https://travis-ci.org/jcmturner/rpc.svg?branch=master)](https://travis-ci.org/jcmturner/rpc)


This project relates to [CDE 1.1: Remote Procedure Call](http://pubs.opengroup.org/onlinepubs/9629399/)

It is a partial implementation that mainly focuses on marshaling NDR encoded byte streams into Go structures.

## Unstable API
Currently this library is at a v0 status to reflect there will be breaking changes in the API without major version revisions.
Please consider this if you adopt this library in your project.

## Help Wanted
* Reference test vectors needed: It has been difficult to implement due to a lack of reference test byte streams in the 
standards documentation. Test driven development has been extremely challenging without these.
If you are aware of and reference test vector sources for NDR encoding please let me know by raising an issue with the details. Thanks!

## References
* [Open Group RPC Publication](http://pubs.opengroup.org/onlinepubs/9629399/)
* [Microsoft RPC Documentation](https://docs.microsoft.com/en-us/windows/desktop/Rpc/rpc-start-page)

## NDR Decode Capability Checklist
- [x] Format label
- [x] Boolean
- [ ] Character
- [x] Unsigned small integer
- [x] Unsigned short integer
- [x] Unsigned long integer
- [x] Unsigned hyper integer
- [ ] Signed small integer
- [ ] Signed short integer
- [ ] Signed long integer
- [ ] Signed hyper integer
- [x] Single float
- [x] Double float
- [x] Uni-dimensional fixed array
- [x] Multi-dimensional fixed array
- [x] Uni-dimensional conformant array
- [ ] Multi-dimensional conformant array
- [x] Uni-dimensional conformant varying array
- [ ] Multi-dimensional conformant varying array
- [ ] Varying string
- [ ] Conformant varying string
- [ ] Array of strings
- [ ] Union
- [ ] Pipe
- [ ] Top level full pointer
- [ ] Top level reference pointer
- [ ] Embedded full pointer
- [ ] Embedded reference pointer