#!/usr/bin/perl -w

use strict;
use warnings;
use POSIX qw(strftime);
use Time::HiRes qw(usleep);

our $runId = 1;

$| = 1;
use sigtrap 'handler' => \&sigtrap, 'INT', 'TERM';

sub sigtrap {
    print "gentext: interrupted, exiting!\n";
    exit(1);
}

sub friendlyDate {
    return strftime "%a %b %e %H:%M:%S %Y", localtime;
}

sub writeText {
    my ( $totalLines, @lines ) = @_;
    my $start = int( rand( $totalLines - 1000 ) );
    my $len   = int( rand(50) ) + 2;

    for ( my $i = $start ; $i < $start + $len ; $i++ ) {
        my $ln = $i - $start + 1;
        print( sprintf( "%3d", $ln ) . ": $lines[$i]\n" );
    }
    print "#### ${\friendlyDate()} #### $len lines printed, starting from $start ($runId). ####\n";
    $runId++;
}

open my $fh, '<', "shakespeare.txt";
chomp( my @lines = <$fh> );
close $fh;
my $linesCount = scalar @lines;

while (1) {
    writeText( $linesCount, @lines );
    usleep( int( rand(1500000) ) + 25000 );
}
