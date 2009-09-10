require 'rubygems'
require 'sinatra'
require 'time'

require 'models'
require 'helpers'

set :static, true

get '/' do
   @vessels = Vessel.all(:order => [:created_at.desc])
   erb :list
end

get '/new' do
   erb :new
end

post '/new' do
   parsed = Vessel.parse(params[:vessel_data])
   if parsed == nil
      @error = 'INVALID VEHICLE CONFIGURATION DATA'
      erb :new
   else
      @vessel = Vessel.first(:cfe=>parsed[:cfe])
      if @vessel == nil
         @vessel = Vessel.new(:data=>parsed[:data], :cfe=>parsed[:cfe])
         if !@vessel.save
             @error = 'ERROR SAVING VESSEL DATA. TRY AGAIN LATER.'
             erb :new
             return
         end
      end
      redirect @vessel.href
   end
end

get '/vessels/?' do
   redirect '/'
end

get '/vessels/:id/?' do
   @vessel = Vessel.get(params[:id])
   erb :show
end
