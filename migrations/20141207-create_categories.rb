class Prec < Sinatra::Base
  migration 'create categories' do
    database.create_table :categories do
      primary_key :id
      column :name, 'varchar(255)', :null => false
    end
  end
end
