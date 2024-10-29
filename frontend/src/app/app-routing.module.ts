import { Component, NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { InicioComponent } from './component/inicio/inicio.component';
import { VisualizadorComponent } from './component/visualizador/visualizador.component';

const routes: Routes = [
  {
    path: 'inicio',
    component: InicioComponent
  },
  { path: 'visualizador',
    component: VisualizadorComponent },
  {
    path: '**',
    redirectTo: 'inicio'
  }
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})
export class AppRoutingModule { }
