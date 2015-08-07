require 'rubygems'
require 'bundler'
Bundler.setup

require 'sinatra'

set :env, :production
disable :run

require_relative './controller'
run Sinatra::Application
