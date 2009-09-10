require 'rubygems'
require 'dm-core'
require 'dm-timestamps'
require 'dm-validations'

DataMapper.setup(:default, {
  :adapter  => 'mysql',
  :database => 'vcd',
  :username => 'root',
  :password => '',
  :host     => 'localhost'
})

class Vessel
   include DataMapper::Resource
   property :id, Serial
   property :created_at, DateTime
   property :cfe, Text
   property :data, Text
   validates_present :cfe, :data

   class <<self
      def parse(data)
         matched = (data =~ /<tt style=\"background-color: rgb\(0,0,0\)\">(.+)<\/tt><br\/><a href=\"http:\/\/www\.captainforever\.com\/captainforever\.php\?cfe=([a-z0-9]+)\">Pilot this vessel<\/a>/m)
         return nil if matched == nil
         {:data=>$1, :cfe=>$2}
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

   def pilot_href
      "http://www.captainforever.com/captainforever.php?cfe=#{cfe}"
   end
end
