# jikei

# install
git clone https://github.com/mkunten/jikei
cd jikei
go -o jikei build
mv jikei </path/to/bin>

# how to use

## jikei list upload
jikei list file: <bid>_coodinate.csv in the dataset
curl -u admin:<password> -X POST http://<server>/jikei/api/admin/jikeilistupload -F "file=@</path/to/csv>"

## page list uplaod
First prepare a page list csv file (identify (ImageMagick) is required)
_tools/sh/pagesize2list.sh <unzipped dataset root dir> > <page_list>.csv 2> error.log
curl -u admin:<password> -X POST http://<server>/jikei/api/admin/pagelistupload -F "file=@</path/to/csv>"
##

## API

parameter | description
-----------------------
<char_id> | character id for IIIF Image API (such like https?://<base_path>/<char_id>/info.json)
<char> | single character
-----------------------

path | description
-----------------------
/api/char/<char_id>/manifest | return IIIF manifest file for character <char_id>
-----------------------
/api/char/search?q=<char> | return JSON file where <char>'s found in the database are listed
-----------------------
