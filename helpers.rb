helpers do
   include Rack::Utils
   alias_method :h, :escape_html

   def eighties_time(time=nil)
      time ||= Time.now
      year = (time.strftime('%Y').to_i - 20).to_s
      year + time.strftime('-%m-%d&nbsp;%H:%M:%S')
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
