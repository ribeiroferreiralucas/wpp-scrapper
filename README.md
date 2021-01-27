# wpp-scrapper

Go module that implements a interface to github.com/Rhymen/go-whatsapp in order to provide simple way to extract all Wpp Messagens and save it to CSV format files.

## Version 1.0 Goals

- [x] Encapsulate authentication methods providing ID and ReAuth
- [x] Method to Start to collect messages
- [x] Store collected messages in a CSV file (one per chat)
- [ ] Collect and store chat infos like name, group or one-to-one, group members (if group), description and others
- [ ] Make program reads config files and env variables to cofigurations like stored file path and other internal configuration thats could be configurable
- [ ] Method to Stop the messages collect
- [ ] WppScrapper should provide a efficient way to API client get each chat scrap stats (running, stoped, finished, queued) and chats list with some util properties (like name, and id).

## Main Problems

- Messages ordenation - The collect messages in each call are unordered if the chuck size are bigger then 1
- Each run collect All messages, even if that are already collected
