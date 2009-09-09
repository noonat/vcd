require 'rubygems'
require 'sinatra'
require 'models'
require 'time'

set :static, true

helpers do
  include Rack::Utils
  alias_method :h, :escape_html
  
  def eighties_time(format, time=nil)
     time ||= Time.now
     time.strftime format
  end
  
  def parse_vcd(data)
     matched = (data =~ /<tt style=\"background-color: rgb\(0,0,0\)\">(.+)<\/tt><br\/><a href=\"http:\/\/www\.captainforever\.com\/captainforever\.php\?cfe=([a-z0-9]+)\">Pilot this vessel<\/a>/m)
     return nil if matched == nil
     {:data=>$1, :cfe=>$2}
  end
  
  def terminal_header(left, opts={})
     s = "<div class=\"#{opts[:classes]}\" style=\"position:relative;\">"
     s << "<div style=\"float:left;\">[#{left}</div>"
     s << "<div style=\"float:right;\">#{opts[:right]}]</div>"
     s << "<div style=\"position:absolute; left:0; text-align:center; width:100%;\">#{opts[:center]}</div>" if opts[:center]
     s << "<div style=\"clear:both;\"></div>"
     s << "</div>"
     s
  end
end

get '/' do
   @vessels = Vessel.all(:order => [:created_at.desc])
   erb :list
end

get '/new' do
   erb :new
end

post '/new' do
   parsed = parse_vcd(params[:vessel_data])
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
