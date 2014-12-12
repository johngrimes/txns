Sequel.migration do
  change do
    create_table :categories do
      primary_key :id
      column :name, 'varchar(255)', :null => false
    end
  end
end
