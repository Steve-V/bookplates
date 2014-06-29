#!/usr/bin/perl -w

use strict;

use Text::CSV;
use WWW::Scraper::ISBN;

my $csv = new Text::CSV;
my @lines = reverse <>;
my $scraper = new WWW::Scraper::ISBN;
$scraper->drivers("LOC", "ISBNnu", "Yahoo");

my $bcid = undef;
my $book = undef;

# Known elephant in Cairo
push @lines, "qr,";

foreach my $line (@lines) {
	$line =~ s/;$//;
	unless ($csv->parse($line)) {
		print "Unparseable", $line, "\n";
		next;
	}
	my ($type, $code) = $csv->fields();
	if ($type eq "qr") {
		if ($book) {
			# upload!
			print "Uploading", $book->{'title'}, "\n";
		}
		$book = undef;
	} else {
		if (!$book) {
			my $r = $scraper->search($code);
			if ($r->found) {
				$book = $r->book;
				print "Found", $code, "\n";
			} else {
				print "Not found", $code, "\n";
			}
		}
	}
}
