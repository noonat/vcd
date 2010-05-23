require 'rubygems'
require 'digest/md5'
require 'dm-core'
require 'dm-timestamps'
require 'dm-types'
require 'dm-validations'
require 'hpricot'

DataMapper.setup(:default, {
  :adapter  => 'mysql',
  :database => 'vcd',
  :username => 'vcd',
  :password => '51k4n2f9',
  :host     => 'localhost'
})

class Vessel
   include DataMapper::Resource
   
   property :id, Serial
   property :created_at, DateTime
   property :ip, IPAddress
   property :cfe, Text, :lazy => false
   property :data, Text, :lazy => false
   
   has n, :vessel_clicks
   has n, :vessel_pilot_clicks
   
   validates_present :cfe, :data
   
   class <<self
      def parse(data)
         matched = (data =~ /<tt style=\"background-color: rgb\(0,0,0\)\">(.+)<\/tt><br\/><a href=\"http:\/\/www\.captainforever\.com\/captainforever\.php\?cfe=([a-z0-9]+)\">Pilot this vessel<\/a>/m)
         return nil if matched == nil
         cfe = $2
         data = Hpricot($1.gsub(/&lt([^;])/, '&lt;\1'))
         data.search('*').each do |node|
            if node.elem?
               case node.name.downcase
               when 'br', 'span'
                  node.attributes.delete_if { |k,v| k.downcase != 'style' }
               else
                  node.parent.children.delete(node)
               end
            elsif node.comment?
               node.parent.children.delete(node)
            end
         end
         data = data.to_s
         return {:cfe=>cfe, :data=>data}
      end
   end
   
   def data_trimmed
      data.split('<br/>').find_all do |line|
         line.gsub(/&nbsp;/, '') != ''
      end.join('<br/>')
   end
   
   def href
      "/vessels/#{id}"
   end
   
   def md5
      return '&lt;null&gt;' if ip.nil?
      Digest::MD5.hexdigest(ip.to_s)
   end

   def pilot_href(track=false)
      return "/vessels/#{id}/pilot" if track
      "http://www.captainforever.com/captainforever.php?cfe=#{cfe}"
   end
end

class VesselClick
   include DataMapper::Resource
   
   property :id, Serial
   property :created_at, DateTime
   property :ip, IPAddress
   property :referrer, URI, :length => 1024
   
   belongs_to :vessel
end

class VesselPilotClick
   include DataMapper::Resource
   
   property :id, Serial
   property :created_at, DateTime
   property :ip, IPAddress
   property :referrer, URI, :length => 1024
   
   belongs_to :vessel
end
