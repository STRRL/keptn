import {ChangeDetectorRef, Component, Input, OnInit, ViewChild} from '@angular/core';
import { DtSort, DtTableDataSource } from '@dynatrace/barista-components/table';

@Component({
  selector: 'ktb-sli-breakdown',
  templateUrl: './ktb-sli-breakdown.component.html',
  styleUrls: ['./ktb-sli-breakdown.component.scss']
})
export class KtbSliBreakdownComponent implements OnInit {

  @ViewChild('sortable', { read: DtSort, static: true }) sortable: DtSort;

  public evaluationState = {
    pass: 'passed',
    warning: 'warning',
    fail: 'failed'
  };

  private _indicatorResults: any;
  private _indicatorResultsFail: any = [];
  private _indicatorResultsWarning: any = [];
  private _indicatorResultsPass: any = [];
  public tableEntries: DtTableDataSource<object> = new DtTableDataSource();

  @Input()
  get indicatorResults(): any {
    return [...this._indicatorResultsFail, ...this._indicatorResultsWarning, ...this._indicatorResultsPass];
  }
  set indicatorResults(indicatorResults: any) {
    if (this._indicatorResults !== indicatorResults) {
      this._indicatorResults = indicatorResults;
      this._indicatorResultsFail = indicatorResults.filter(i => i.status === 'fail');
      this._indicatorResultsWarning = indicatorResults.filter(i => i.status === 'warning');
      this._indicatorResultsPass = indicatorResults.filter(i => i.status !== 'fail' && i.status !== 'warning');
      this.updateDataSource();
      this._changeDetectorRef.markForCheck();
    }
  }

  constructor(private _changeDetectorRef: ChangeDetectorRef) {
  }

  ngOnInit(): void {
    this.sortable.sort('score', 'asc');
    this.tableEntries.sort = this.sortable;
  }

  private updateDataSource() {
    this.tableEntries.data = this.assembleTablesEntries(this.indicatorResults);
  }

  private formatNumber(value: number) {
    let n = value;
    if (n < 1) {
      n = Math.floor(n * 1000) / 1000;
    } else if (n < 100) {
      n = Math.floor(n * 100) / 100;
    } else if (n < 1000) {
      n = Math.floor(n * 10) / 10;
    } else {
      n = Math.floor(n);
    }

    return n;
  }

  private assembleTablesEntries(indicatorResults): any {
    const totalscore  = indicatorResults.reduce((acc, result) => acc + result.score, 0);
    return indicatorResults.map(indicatorResult =>  {
      return {
        name: indicatorResult.value.metric,
        value: this.formatNumber(indicatorResult.value.value),
        result: indicatorResult.status,
        score: this.round(indicatorResult.score / totalscore, 2),
        targets: indicatorResult.targets
      };
    });
  }

  private round(value: number, places: number): number {
    return +(Math.round(Number(`${value}e+${places}`))  + `e-${places}`);
  }

}
