load 'deploy' if respond_to?(:namespace) # cap2 differentiator

default_run_options[:pty] = true

set :application, 'vcd'
set :user, 'webuser'
set :use_sudo, false

set :scm, :git
set :scm_verbose, true
set :repository, "git@github.com:noonat/vcd.git"
set :deploy_via, :remote_cache
set :deploy_to, "/home/#{user}/#{application}"
set :branch, 'master'

role :app, "vcd.phuce.com"
role :web, "vcd.phuce.com"
role :db, "vcd.phuce.com", :primary => true

set :runner, user
set :admin_runner, user

namespace :deploy do
   task :restart do
      run "touch #{current_path}/tmp/restart.txt"
   end
end
