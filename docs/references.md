# elbow

- [elbow](#elbow)
  - [Project home](#project-home)
  - [Packages](#packages)
  - [Inspiration, Guidance](#inspiration-guidance)
    - [Configuration object](#configuration-object)
    - [Sorting files](#sorting-files)
    - [Path/File Existence](#pathfile-existence)
    - [Slice management](#slice-management)

## Project home

See [our GitHub repo](https://github.com/atc0005/elbow) for the latest code,
to file an issue or submit improvements for review and potential inclusion
into the project.

## Packages

The following list of packages were either used or seriously considered at
some point during the development of this project.

- <https://github.com/jessevdk/go-flags#example>
- <https://github.com/sirupsen/logrus>
- <https://github.com/integrii/flaggy>

## Inspiration, Guidance

The following unordered list of sites/examples provided guidance while
developing this application. Depending on when consulted, the original code
written based on that guidance may no longer be present in the active version
of this application.

### Configuration object

- <https://github.com/go-sql-driver/mysql/blob/877a9775f06853f611fb2d4e817d92479242d1cd/dsn.go#L67>
- <https://github.com/aws/aws-sdk-go/blob/10878ad0389c5b3069815112ce888b191c8cd325/aws/config.go#L251>
- <https://github.com/aws/aws-sdk-go/blob/master/aws/config.go>
- <https://github.com/aws/aws-sdk-go/blob/10878ad0389c5b3069815112ce888b191c8cd325/awstesting/integration/performance/s3GetObject/config.go#L25>
- <https://github.com/aws/aws-sdk-go/blob/10878ad0389c5b3069815112ce888b191c8cd325/awstesting/integration/performance/s3GetObject/main.go#L25>

### Sorting files

- <https://stackoverflow.com/questions/46746862/list-files-in-a-directory-sorted-by-creation-time>

### Path/File Existence

- <https://gist.github.com/mattes/d13e273314c3b3ade33f>

### Slice management

- <https://yourbasic.org/golang/delete-element-slice/>
- <https://stackoverflow.com/questions/37334119/how-to-delete-an-element-from-a-slice-in-golang>
- <https://github.com/golang/go/wiki/SliceTricks>
