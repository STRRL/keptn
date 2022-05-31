/// <reference types="cypress" />

import { Project } from '../../../client/app/_models/project';
import { interceptDashboard } from '../intercept';
import EnvironmentPage from './EnvironmentPage';
import NewProjectCreatePage from './NewProjectCreatePage';

class DashboardPage {
  private PROJECT_TILE_LOCATOR = 'dt-tile[id="proj_pattern"]';
  private STAGE_TAG_LOCATOR = 'dt';

  public intercept(): this {
    interceptDashboard();
    return this;
  }

  public visit(): this {
    cy.visit(`/`).wait('@metadata');
    return this;
  }

  public clickProjectTile(projectName: string): EnvironmentPage {
    cy.wait('@projects').get(this.PROJECT_TILE_LOCATOR.replace('proj_pattern', projectName)).click();
    return new EnvironmentPage();
  }

  public assertProjects(projects: Project[]): this {
    cy.get('ktb-project-tile').should('have.length', projects.length);
    projects.forEach((project, index) => {
      cy.get('ktb-project-tile').eq(index).find('dt-tile-title').should('contain.text', project.projectName);
      cy.get('ktb-project-tile')
        .eq(index)
        .byTestId('keptn-project-tile-numStagesServices')
        .should('contain.text', `${project.stages.length} Stages, ${project.stages[0].services.length} Services `);
    });
    return this;
  }

  public clickCreateNewProjectButton(): NewProjectCreatePage {
    cy.get('.dt-button-primary > span.dt-button-label').contains('Create a new project').click();
    return new NewProjectCreatePage();
  }

  public clickStageTag(projectName: string, stageName: string): EnvironmentPage {
    cy.byTestId(`ktb-project-${projectName}`).find('dt-tag a').contains(stageName).click();
    return new EnvironmentPage();
  }
}

export default DashboardPage;
