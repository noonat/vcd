namespace :db do
    task :upgrade do
        require 'models'
        DataMapper.auto_upgrade!
    end
end
