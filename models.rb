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

   def href
      "/vessels/#{id}"
   end

   def pilot_href
      "http://www.captainforever.com/captainforever.php?cfe=#{cfe}"
   end
end
