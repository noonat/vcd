namespace :db do
    task :migrate do
        require 'models'
        require 'dm-migrations'
        DataMapper.auto_upgrade!
    end
end
