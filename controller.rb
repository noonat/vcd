require 'rubygems'
require 'sinatra'
require 'time'
require 'haml'

require_relative './models'
require_relative './helpers'

set :static, true

get '/' do
   @vessels = Vessel.all(:order=>[:created_at.desc]) #, :limit=>30)
   rows = repository(:default).adapter.query(
      'SELECT vessel_id, COUNT(*) AS count
       FROM vessel_pilot_clicks GROUP BY vessel_id')
   @vessel_pilot_clicks = rows.inject({}) do |clicks, row|
      clicks[row.vessel_id] = row.count
      clicks
   end
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
         @vessel = Vessel.new(:data=>parsed[:data], :cfe=>parsed[:cfe], :ip=>request.ip)
         if !@vessel.save
             @error = 'ERROR SAVING VESSEL DATA. TRY AGAIN LATER.'
             erb :new
             return
         end
      else
         click = @vessel.vessel_pilot_clicks.build(:referrer=>request.referrer, :ip=>request.ip)
         click.save
      end
      redirect @vessel.href
   end
end

get '/rss.xml' do
   @vessels = Vessel.all(:order=>[:created_at.desc], :limit=>30)
   haml :rss
end

get '/vessels/?' do
   redirect '/'
end

get '/vessels/:id/?' do
   @vessel = Vessel.get(params[:id])
   raise Sinatra::NotFound if @vessel.nil?
   click = @vessel.vessel_clicks.build(:referrer=>request.referrer, :ip=>request.ip)
   click.save
   erb :show
end

get '/vessels/:id/pilot/?' do
   vessel = Vessel.get(params[:id])
   raise Sinatra::NotFound if vessel.nil?
   click = vessel.vessel_pilot_clicks.build(:referrer=>request.referrer, :ip=>request.ip)
   click.save
   redirect vessel.pilot_href
end
