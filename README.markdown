# Vehicle Configuration Data Archive

Allows users to paste their [Captain Forever](http://captainforever.com) vessel exports to share with others.

## Installation

vcd requires:

- [DataMapper](http://datamapper.org) to talk MySQL
- [Sinatra](http://sinatrarb.com) to talk HTTP
- [Haml](http://haml-lang.com/) to talk HTML
- [Hpricot](http://github.com/whymirror/hpricot/tree/master) to sanitize things

You can install them with gem:

    sudo gem install datamapper do_mysql haml hpricot sinatra

## Usage

You can run the app locally using Sinatra:

    $ ruby controller.rb 
    == Sinatra/0.9.4 has taken the stage on 4567 for development with backup from Mongrel

vcd also has a config.ru file for use with Passenger.
