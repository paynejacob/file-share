# File Share

---

File Share is a dead simple application for moving files from one machine to another.  The first page you see will ask for a file, the next will give you a short link to download that file.  That is it. After 15 minutes the link expires, and the file is deleted.  There is no need to cleanup old files, login, or install any software.  Just put the file on the site and download it with the link provided.  File Share was created when I became frustrated trying to take a picture of a document on my phone to upload to my work computer.

## Installation

File share can be run as a standalone docker container.

```shell
docker run -p 80:80 paynejacob/file-share:latest
```

It can also be installed as a helm chart

```shell
git clone https://github.com/file-share.git
helm upgrade --install --create-namespace -n file-share file-share file-share/charts/file-share
```

## Usage

### In your browser

The simples way to use file share is with a browser.  Simply navigate to your url and follow the on screen instructions.

### With curl

If you want to use file share from the command line you can use curl.

```shell
curl -Ls -w "%{url_effective}?download\n" -o /dev/null -F file=@<local file> <file share url>
```

## Contributing

Have an idea for a feature? Found a bug? Please [create an issue.](https://github.com/paynejacob/file-share/issues/new)

Pull requests are always welcome!  If you want to resolve an issue, please make sure it is not assigned to anyone before starting on it.  The assignee's pr will always be given favor.

## License

Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with the License. You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.
