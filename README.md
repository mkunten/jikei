# jikei

## prepare db

```bash
createuser -P jikei
# set <password>
createdb -O jikei jikei
```

## build

```bash
git clone https://github.com/mkunten/jikei
cd jikei
go get
swag init
vi config.go
go build -o jikei 
# or go build -o jikei -ldflags "-X main.version=<version>"
mv jikei </path/to/bin>
```

## run

```bash
</path/to/bin>/jikei
```

## prepare jikei list files

Jikei list files are CSV files named <bid>_coordinate.csv. Download the ZIP files from the CODH dataset (https://codh.rois.ac.jp/char-shape/) and unzip them to a directory (e.g., ./data).

## register data

### register jikei lists

```bash
export JIKEI_ADMIN_PASSWORD="password"
# register data per csv file
curl -u admin:$JIKEI_ADMIN_PASSWORD -X POST http://<server>/jikei/api/admin/jikeilistupload -F "file=@</path/to/csv>"
```

### register a page list
First prepare a page list csv file (`identify` (ImageMagick) is required)

```bash
_tools/sh/pagesize2list.sh <unzipped dataset root dir> > <page_list>.csv 2> error.log
curl -u admin:$JIKEI_ADMIN_PASSWORD -X POST http://<server>/jikei/api/admin/pagelistupload -F "file=@</path/to/csv>"
```

## API

|parameter | description|
|----------|------------|
|&lt;char_id&gt; | character id for IIIF Image API (such like https?://&lt;base_path&gt;/&lt;char_id&gt;/info.json)|
|&lt;char&gt; | single character|

|method |path | description|
|-------|-----|-----------------|
|GET|/jikei/api/char/&lt;char_id&gt;/manifest | return IIIF manifest file for character &lt;char_id&gt;|
|GET|/jikei/api/char/search?q=&lt;char&gt; | return JSON file where &lt;char&gt;&apos;s found in the database are listed|

## API docs (partial)

/jikei/api/docs
