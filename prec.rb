require 'sinatra/base'
require 'sinatra/sequel'
require 'json'
require 'csv'

require 'pry'

class Prec < Sinatra::Base
  register Sinatra::SequelExtension

  Dir[
    File.join(File.dirname(__FILE__), 'migrations', '*.rb')
  ].each { |file| require file}

  get '/txns.html' do
    @txns = database[:txns].reverse_order(:date)
    @categories = database[:categories].order(:name)
    erb :txns
  end

  post '/txns' do
    tempfile = params[:txn_file][:tempfile]
    txns = []
    CSV.foreach(tempfile, :headers => :first_row) do |row|
      hash = hash_txn_row(row.to_a.map { |x| x.last })
      next unless database[:txns].where(:hash => hash).empty?
      date = Time.parse(row[0]).iso8601
      description = row[1]
      debit_cents = parse_currency_string(row[2], true)
      credit_cents = parse_currency_string(row[3], true)
      balance_cents = parse_currency_string(row[4])
      database[:txns].insert(:hash => hash, :date => date,
        :description => description, :debit_cents => debit_cents,
        :credit_cents => credit_cents, :balance_cents => balance_cents)
    end
    redirect to('/txns.html')
  end

  def hash_txn_row(row)
    concat_row = row.join('')
    hash = Digest::SHA2.new(512)
    hash.update(concat_row)
    hash.hexdigest
  end

  def parse_currency_string(string, force_positive = false)
    return nil unless string
    value = string.gsub!('.', '').to_i
    if force_positive and value < 0
      value * -1
    else
      value
    end
  end

  def format_currency(cents)
    return '0.00' if cents.nil?
    dollars = cents / 100.0
    "%.2f" % dollars
  end
end
