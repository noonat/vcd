require 'rubygems'
require 'sinatra'

set :env, :production
disable :run

require 'controller'
run Sinatra::Application
