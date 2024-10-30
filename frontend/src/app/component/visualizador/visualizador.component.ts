import { Component, OnInit } from '@angular/core';
import { ApiService } from '../../services/api.service';
import {FormsModule} from "@angular/forms";
import {NgForOf, NgStyle} from "@angular/common";

@Component({
  selector: 'app-visualizador',
  templateUrl: './visualizador.component.html',
  standalone: true,
  imports: [
    FormsModule,
    NgForOf,
    NgStyle
  ],
  styleUrls: ['./visualizador.component.css']
})
export class VisualizadorComponent implements OnInit {
  // Propiedad para almacenar la respuesta de la API
  entrada = "";
  data: any[] = [];
  selectedItem: any;
  salida: string = "";

  constructor(public service: ApiService) { }

  ngOnInit(): void {
    // Llama a la API con un 'path' predeterminado
    this.entrada = this.entrada || '/home/ubuntu/';
    this.buscar(); // Realiza la búsqueda inicial
  }

  buscar(): void {
    if (this.entrada.trim() !== "") {
      this.service.getDatos(this.entrada).subscribe((res: any[]) => {
        this.data = res;
      }, (error: any) => {
        console.error("Error al enviar el comando:", error);
      });
    } else {
      alert("Cadena de entrada vacía...");
    }
  }


  selectItem(item: any) {
    if (!item.is_dir && item.name.endsWith('.txt')) {
      this.selectedItem = item;
    } else {
      this.selectedItem = null;
    }
  }

  viewContent() {
    if (this.selectedItem) {
      this.service.getFileContent(this.selectedItem.path).subscribe((content: string) => {
        this.salida = content; // Asigna el contenido del archivo al área de texto
      }, (error: any) => {
        console.error("Error al obtener el contenido del archivo:", error);
        this.salida = "Error al obtener el contenido del archivo."; // Mensaje de error
      });
    }
  }


}
