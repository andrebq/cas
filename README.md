# What is this?

CAS stands for Content-Adressable Storage. It is a special case of storage where
you access the content by its cryptographic hash.

There is no update operation, in order to change a content you need to "insert" it
again.

# License?

ISC License (aka MIT). Read LICENSE file for more details

# Status

For now, it is a WIP. It will be used as part of another more interesting project.

# Can I use it as a stand-alone storage?

Not now but in the future! I'll add a gRPC and HTTP/2 API's on top of it with 
the option of using Redis as main-storage or cache-layer. Since data is immutable, caching is simple.

# What about delete operations?

Right now there is no way to delete a hash (ofcourse one can open the backend and delete it there). In the future MAYBE some form of garbage collection can be
implemented (maybe using atomic-ref-count).

The main issue is, once one assumes a hash exists we there is no external management of "how many systems reference this specific hash" it is hard to delete "right away".

For cases where data MUST BE erased (right to be forgotten or legal requirements).
I'll design an alternative where a SHA1 hash "updated" with a "corrupted" content, but that content MUST BE signed with a "trusted" key. This, I think, is a good
trade-off between the system's reliability/security and the "changing nature" of our world.

For now "trusted", "updated", "corrupted" are too vague to be defined at this point.