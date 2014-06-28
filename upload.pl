#!/usr/bin/perl -w

use strict;

use Text::CSV;

my $csv = new Text::CSV;
my @lines = reverse <>;

foreach my $line (@lines) {
	$line =~ s/;$//;
	print $line;
	print join(";", $csv->fields()), "\n" if $csv->parse($line);
}
