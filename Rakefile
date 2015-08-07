namespace :db do
    task :migrate do
        require_relative './models'
        require 'dm-migrations'
        DataMapper.auto_upgrade!
    end
end
