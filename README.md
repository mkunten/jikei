# jikei

# install
git clone https://github.com/mkunten/jikei  
cd jikei  
go get  
go -o jikei build  
mv jikei &lt;/path/to/bin&gt;

# how to use

## jikei list upload
jikei list file: &lt;bid&gt;_coodinate.csv in the dataset

```bash
curl -u admin:<password> -X POST http://<server>/jikei/api/admin/jikeilistupload -F "file=@</path/to/csv>"
```

## page list upload
First prepare a page list csv file (identify (ImageMagick) is required)

```bash
_tools/sh/pagesize2list.sh <unzipped dataset root dir> > <page_list>.csv 2> error.log
curl -u admin:<password> -X POST http://<server>/jikei/api/admin/pagelistupload -F "file=@</path/to/csv>"
```

## API

|parameter | description|
|----------|------------|
|&lt;char_id&gt; | character id for IIIF Image API (such like https?://&lt;base_path&gt;/&lt;char_id&gt;/info.json)|
|&lt;char&gt; | single character|

|method |path | description|
|-------|-----|-----------------|
|GET|/api/char/&lt;char_id&gt;/manifest | return IIIF manifest file for character &lt;char_id&gt;|
|GET|/api/char/search?q=&lt;char&gt; | return JSON file where &lt;char&gt;&apos;s found in the database are listed|
