## rangastalam

rangastalam is an experimental command line video editor being written in go

it uses a simple [DSL](https://wikipedia.org/wiki/Domain-specific_language) to describe timelines of video, audio, text, and images,  
and then translates that into [ffmpeg](https://ffmpeg.org) commands  

the long term goal is to make it easy to create video essays in a scriptable way :
- cut and arrange multiple clips
- add audio tracks and voiceovers
- overlay text, images, or additional videos
- apply transforms ( crop, scale, rotate, move ) and animations
- render everything reproducibly from a single script

right now this project is in early development, expect things to change a lot

