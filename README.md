# Moodle File Renamer

![Static Badge](https://img.shields.io/badge/Version-1.0.0-blue)

![Static Badge](https://img.shields.io/badge/License-MIT_License-orange)

## :speaking_head: Overview

Moodleから一括ダウンロードしたファイル群は，ユーザごとにフォルダが作成され，そのフォルダに提出されたファイルが置かれる．そのユーザごとのフォルダは余計な情報が多く長くなるため，リネームして利用することが多い．

ただ，リネームするスクリプトをその都度作成するのも面倒だし，同じ階層に異なるファイルやディレクトリが置かれていた場合，そのファイル，ディレクトリもリネームの対象となる．
加えて，元のファイル名に戻してやり直したいという要求も場合によってはあり得るであろう．

そこで，一括して所定の名前にリネームでき，元のファイル名に戻すことも可能なツールを作成した．

## How to use

```sh
mfr version 1.0.0
Usage: mfr [OPTIONS] <DIR>
OPTIONS
  -r, --restoration      restores the renamed directory names to the original.
  -f, --format <FORMAT>  specifies the format of the resultant directory names.
                         default: "%default" (is equals to "%uid_%lname_%fname")
                         available variables: default, original, uid, fname,
                         lname, name, sid, and note.
  -h, --help             prints this help message and exit.
ARGUMENTS
  DIR    the target directory containing the downloaded directories from Moodle.
```

