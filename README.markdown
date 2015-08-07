# Vehicle Configuration Data Archive

Allows users to paste their [Captain Forever](http://captainforever.com) vessel exports to share with others.

## Installation

vcd requires:

- [DataMapper](http://datamapper.org) to talk MySQL
- [Sinatra](http://sinatrarb.com) to talk HTTP
- [Haml](http://haml-lang.com/) to talk HTML
- [Hpricot](http://github.com/whymirror/hpricot/tree/master) to sanitize things

You can install them with bundler:

    bundle install

## Usage

You can run the app locally using Rack:

    $ bundle exec rackup

You can also run it using Docker:

    $ docker build -t noonat/vcd .
    $ docker run --rm --name vcd --publish 8000:80 --env VCD_DB_HOST=172.17.42.1 noonat/vcd
