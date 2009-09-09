require 'rubygems'
require 'sinatra'

Sinatra::Application.default_options.merge!(
	:run => false,
	:env => :production
)

require 'controller'
run Sinatra.application
