require 'rubygems'
require 'bundler'
Bundler.setup

require 'sinatra'

set :env, :production
disable :run

require 'controller'
run Sinatra::Application
