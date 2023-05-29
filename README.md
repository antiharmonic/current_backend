# current - a project to help track media consumption

A standalone python version exists at https://github.com/antiharmonic/current_standalone

I wanted to play around with making a backend in the cloud so I could use this functionality wherever I went, in addition to getting better at Go and software architecture in general.

As such this Go backend separates out the transport and storage layers, using a pattern I found by @morganhein via https://github.com/morganhein/backend-takehome-telegraph. I'd also like to add rabbit in here somewhere, but I'm not entirely sure on how. Perhaps this monolith can be a consumer too, or do I write a separate service to do that?
