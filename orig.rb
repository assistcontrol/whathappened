#!/usr/local/bin/ruby --yjit

# Make sure we speak the same encoding as git
Encoding.default_external = 'UTF-8'
Encoding.default_internal = 'UTF-8'

require 'date'
require 'optparse'
require 'tempfile'

@Root     = '/data/freebsd'
Repos     = %w[ports src doc]
PortsRepo = 'ports'
OSVer     = %x[/usr/bin/uname -U].to_i / 10**5  # 180200 -> 18
LogFmt = {
  ports: '━' * 20 + %q[%nCommitter: %cl (%cn)%nDate: %cd%nCommit: https://cgit.freebsd.org/ports/commit/?id=%h%n%n%B],
  src:   '━' * 20 + %q[%nCommitter: %cl (%cn)%nDate: %cd%nCommit: https://cgit.freebsd.org/src/commit/?id=%h%n%n%B],
}
DateFmt = '%a %d %b %H:%M'

@opts = {
  date:    Date.today - 1,
  fast:    false,
  verbose: false
}
OptionParser.new do |o|
  o.on('-dDATE', '--date=DATE', 'Run for specified date (YYYY-MM-DD)') {|d| @opts[:date] = Date.strptime(d, '%Y-%m-%d') }
  o.on('-v', '--verbose', 'Show progress') { @opts[:verbose] = true }
  o.on('-h', '--help',    'Shows this help') { puts o; exit }
end.parse!

def debug(s)
  puts s if @opts[:verbose]
end

# Only show each commit once
@seen = {}
def filter(commits)
  good = commits.flatten.sort.uniq.delete_if {|h| @seen[h] }
  good.map {|c| @seen[c] = 1; c.split(' ')[1] }.join(' ')
end

# Show a formatted log message for a given list of commits
def logs(repo, commits, format)
  commitlist = filter(commits)
  cmd = %Q[/usr/local/bin/git -C "#{@Root}/#{repo}" show --no-patch --date=format-local:'#{DateFmt}' --format='#{format}' #{commitlist}]
  debug "logs: ``#{cmd}''"
  %x[#{cmd}]
end

# Obtain a list of revisions that match a set of queries
def revlist(repo, limiters)
  cmd = %Q[/usr/local/bin/git -C "#{@Root}/#{repo}" log #{@dateRange} --format='%ct %H' #{limiters}]
  debug "rev list: ``#{cmd}''"
  %x[#{cmd}].split("\n")
end

def title(s)
  "███ #{s}\n\n"
end

#
# DO STUFF:
@dateRange = "--since #{@opts[:date].strftime('%Y-%m-%d')}:00:00 --before #{(@opts[:date] + 1).strftime('%Y-%m-%d')}:00:00"

# Update Repos
debug title('prepping repos')
Repos.each do |repo|
  debug "-----<<(#{repo})>>-----"

  cmd = %Q(/usr/local/bin/git -C "#{@Root}/#{repo}" pull -q)
  debug cmd
  %x[#{cmd}]
end

## Get list of local ports
ports = %x|/usr/local/sbin/pkg search -o '.*'|.split("\n")
        .collect {|s| s.split[0] }
        .keep_if {|d| Dir.exist? "#{@Root}/#{PortsRepo}/#{d}"}

## Get list of my ports
index = Dir['/usr/ports/INDEX-1[0-9]'].sort.reverse.first
File.open(index).grep('adamw@FreeBSD.org').each do |l|
  parts = l.split('|')
  ports << parts[1].delete_prefix('/usr/ports/')
end
ports.sort.uniq!

# Relevant commits
revs =  revlist('ports', '--committer adamw@FreeBSD.org')
revs << revlist('ports', '--grep adamw')
revs << revlist('ports', ports.join(' '))
relevant = logs('ports', revs, LogFmt[:ports])

# Other ports
revs  = revlist('ports', 'Mk Tools Templates')
other = logs('ports', revs, LogFmt[:ports])

# src
revs = revlist('src', "stable/#{OSVer}")
src  = logs('src', revs, LogFmt[:src])

unless relevant.empty?
  puts title('relevant ports')
  puts relevant
end

unless other.empty?
  puts title('other ports')
  puts other
end

unless src.empty?
  puts title('src')
  puts src
end
