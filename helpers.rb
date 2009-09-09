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

   def tty_dialog(opts={}, &block)
      opts[:width]  ||= '812px'
      opts[:height] ||= '96px'
      opts[:top]    ||= '28px'
      opts[:block]  ||= block
      @_out_buf << erb(:_tty_dialog, :layout => false, :locals => opts)
   end
   
   def tty_h1(text, opts={})
      s = <<-HTML
      <div class="tty_h1 #{opts[:classes]}">
         <div class="left">[#{text}</div>
         <div class="right">#{opts[:right]}]</div>
      HTML
      s << "<div class=\"center\">#{opts[:center]}</div>" if opts[:center]
      s << <<-HTML
         <div style="clear:both;"></div>
      </div>
      HTML
      s
   end
   
   def tty_h2(text)
      <<-HTML
      <div class="tty_h2"><span>#{text}</span></div>
      HTML
   end
end
