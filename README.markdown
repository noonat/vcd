# Vehicle Configuration Data Archive

Allows users to paste their [Captain Forever](http://captainforever.com) vessel exports to share with others.

## Installation

VCD requires [DataMapper](http://datamapper.org) and [Sinatra](http://sinatrarb.com).

You can install them with gem:

    sudo gem install datamapper do_mysql sinatra

## Usage

You can run the app locally using Sinatra:

    $ ruby controller.rb 
    == Sinatra/0.9.4 has taken the stage on 4567 for development with backup from Mongrel

VCD also has a config.ru file for use with Passenger.
