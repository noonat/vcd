require 'rubygems'
require 'sinatra'

Sinatra::Default.set(:run, false)
Sinatra::Default.set(:env, :production)

require 'controller'
run Sinatra::Application
