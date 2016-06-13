Sequel.migration do
  change do
    create_table :imports do
      primary_key :id
      column :hashes, 'text[]', :null => false
      timestamp :datetime, :null => false
    end
  end
end
